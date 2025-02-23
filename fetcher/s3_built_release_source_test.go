package fetcher_test

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/kiln/fetcher"
	"github.com/pivotal-cf/kiln/fetcher/fakes"
	"github.com/pivotal-cf/kiln/internal/cargo"
)

var _ = Describe("GetMatchedReleases from S3 built source", func() {
	var (
		releaseSource     fetcher.S3BuiltReleaseSource
		fakeS3Client      *fakes.S3ObjectLister
		desiredReleaseSet fetcher.ReleaseSet
		bpmKey            string
		ignoredStemcell   cargo.Stemcell
	)

	BeforeEach(func() {
		bpmReleaseID := fetcher.ReleaseID{Name: "bpm", Version: "1.2.3-lts"}
		desiredReleaseSet = fetcher.ReleaseSet{
			bpmReleaseID: fetcher.CompiledRelease{
				ID:              bpmReleaseID,
				StemcellOS:      "ubuntu-xenial",
				StemcellVersion: "190.0.0",
				Path:            "",
			},
		}

		fakeS3Client = new(fakes.S3ObjectLister)

		irrelevantKey := "some-key"
		uaaKey := "1.10/uaa/uaa-1.2.3.tgz"
		bpmKey = "2.5/bpm/bpm-1.2.3-lts.tgz"
		fakeS3Client.ListObjectsPagesStub = func(input *s3.ListObjectsInput, fn func(*s3.ListObjectsOutput, bool) bool) error {
			shouldContinue := fn(&s3.ListObjectsOutput{
				Contents: []*s3.Object{
					{Key: &irrelevantKey},
					{Key: &uaaKey},
					{Key: &bpmKey},
				},
			},
				true,
			)
			Expect(shouldContinue).To(BeTrue())
			return nil
		}

		ignoredStemcell = cargo.Stemcell{OS: "ignored", Version: "ignored"}

		logger := log.New(nil, "", 0)

		releaseSource = fetcher.S3BuiltReleaseSource{
			Logger:   logger,
			S3Client: fakeS3Client,
			Regex:    `^2.5/.+/(?P<release_name>[a-z-_]+)-(?P<release_version>[0-9\.]+(-\w+(\.[0-9]+)?)?)\.tgz$`,
			Bucket:   "built-bucket",
		}
	})

	It("lists all objects that match the given regex", func() {
		matchedS3Objects, err := releaseSource.GetMatchedReleases(desiredReleaseSet, ignoredStemcell)
		Expect(err).NotTo(HaveOccurred())

		input, _ := fakeS3Client.ListObjectsPagesArgsForCall(0)
		Expect(input.Bucket).To(Equal(aws.String("built-bucket")))

		Expect(matchedS3Objects).To(HaveLen(1))
		Expect(matchedS3Objects).To(HaveKeyWithValue(fetcher.ReleaseID{Name: "bpm", Version: "1.2.3-lts"}, fetcher.BuiltRelease{ID: fetcher.ReleaseID{Name: "bpm", Version: "1.2.3-lts"}, Path: bpmKey}))
	})

	When("the regular expression is missing a capture group", func() {
		BeforeEach(func() {
			releaseSource = fetcher.S3BuiltReleaseSource{
				Regex: `^2.5/.+/([a-z-_]+)-(?P<release_version>[0-9\.]+)\.tgz$`,
			}
		})

		It("returns an error if a required capture is missing", func() {
			_, err := releaseSource.GetMatchedReleases(nil, ignoredStemcell)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring("Missing some capture group")))
		})
	})

	Context("if any objects in built S3 bucket do not match a release specified in Kilnfile.lock", func() {
		BeforeEach(func() {
			wrongReleaseVersionKey := "2.5/bpm/bpm-4.5.6.tgz"
			wrongReleaseNameKey := "2.5/diego/diego-1.2.3.tgz"
			fakeS3Client.ListObjectsPagesStub = func(input *s3.ListObjectsInput, fn func(*s3.ListObjectsOutput, bool) bool) error {
				shouldContinue := fn(&s3.ListObjectsOutput{
					Contents: []*s3.Object{
						{Key: &wrongReleaseVersionKey},
						{Key: &wrongReleaseNameKey},
						{Key: &bpmKey},
					},
				},
					true,
				)
				Expect(shouldContinue).To(BeTrue())
				return nil
			}
		})

		It("does not return them, but does return the matched release", func() {
			matchedS3Objects, err := releaseSource.GetMatchedReleases(desiredReleaseSet, ignoredStemcell)
			Expect(err).NotTo(HaveOccurred())

			Expect(matchedS3Objects).To(HaveLen(1))
			Expect(matchedS3Objects).To(HaveKeyWithValue(fetcher.ReleaseID{Name: "bpm", Version: "1.2.3-lts"}, fetcher.BuiltRelease{ID: fetcher.ReleaseID{Name: "bpm", Version: "1.2.3-lts"}, Path: bpmKey}))
		})
	})
})

var _ = Describe("S3BuiltReleaseSource DownloadReleases from Built source", func() {
	var (
		logger           *log.Logger
		releaseSource    fetcher.S3BuiltReleaseSource
		releaseDir       string
		matchedS3Objects map[fetcher.ReleaseID]fetcher.ReleaseInfoDownloader
		fakeS3Downloader *fakes.S3Downloader
	)

	BeforeEach(func() {
		var err error

		releaseDir, err = ioutil.TempDir("", "kiln-releaseSource-test")
		Expect(err).NotTo(HaveOccurred())

		matchedS3Objects = make(map[fetcher.ReleaseID]fetcher.ReleaseInfoDownloader)
		matchedS3Objects[fetcher.ReleaseID{Name: "uaa", Version: "1.2.3"}] = fetcher.BuiltRelease{ID: fetcher.ReleaseID{Name: "uaa", Version: "1.2.3"}, Path: "some-uaa-key"}
		matchedS3Objects[fetcher.ReleaseID{Name: "bpm", Version: "1.2.3"}] = fetcher.BuiltRelease{ID: fetcher.ReleaseID{Name: "bpm", Version: "1.2.3"}, Path: "some-bpm-key"}

		logger = log.New(GinkgoWriter, "", 0)
		fakeS3Downloader = new(fakes.S3Downloader)
		// fakeS3Downloader writes the given S3 bucket and key into the output file for easy verification
		fakeS3Downloader.DownloadStub = func(writer io.WriterAt, objectInput *s3.GetObjectInput, setConcurrency ...func(dl *s3manager.Downloader)) (int64, error) {
			n, err := writer.WriteAt([]byte(fmt.Sprintf("%s/%s", *objectInput.Bucket, *objectInput.Key)), 0)
			return int64(n), err
		}
		releaseSource = fetcher.S3BuiltReleaseSource{
			Logger:       logger,
			S3Downloader: fakeS3Downloader,
			Bucket:       "some-bucket",
		}
	})

	AfterEach(func() {
		_ = os.RemoveAll(releaseDir)
	})

	It("downloads the appropriate versions of built releases listed in matchedS3Objects", func() {
		err := releaseSource.DownloadReleases(releaseDir, matchedS3Objects, 7)
		Expect(err).NotTo(HaveOccurred())
		Expect(fakeS3Downloader.DownloadCallCount()).To(Equal(2))

		bpmContents, err := ioutil.ReadFile(filepath.Join(releaseDir, "bpm-1.2.3.tgz"))
		Expect(err).NotTo(HaveOccurred())
		Expect(bpmContents).To(Equal([]byte("some-bucket/some-bpm-key")))
		uaaContents, err := ioutil.ReadFile(filepath.Join(releaseDir, "uaa-1.2.3.tgz"))
		Expect(err).NotTo(HaveOccurred())
		Expect(uaaContents).To(Equal([]byte("some-bucket/some-uaa-key")))

		_, _, opts := fakeS3Downloader.DownloadArgsForCall(0)
		verifySetsConcurrency(opts, 7)

		_, _, opts = fakeS3Downloader.DownloadArgsForCall(1)
		verifySetsConcurrency(opts, 7)
	})

	Context("when the matchedS3Objects argument is empty", func() {
		It("does not download anything from S3", func() {
			err := releaseSource.DownloadReleases(releaseDir, fetcher.ReleaseSet{}, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeS3Downloader.DownloadCallCount()).To(Equal(0))
		})
	})

	Context("when number of threads is not specified", func() {
		It("uses the s3manager package's default download concurrency", func() {
			err := releaseSource.DownloadReleases(releaseDir, matchedS3Objects, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeS3Downloader.DownloadCallCount()).To(Equal(2))

			_, _, opts1 := fakeS3Downloader.DownloadArgsForCall(0)
			verifySetsConcurrency(opts1, s3manager.DefaultDownloadConcurrency)

			_, _, opts2 := fakeS3Downloader.DownloadArgsForCall(1)
			verifySetsConcurrency(opts2, s3manager.DefaultDownloadConcurrency)
		})
	})

	Context("failure cases", func() {
		Context("when a file can't be created", func() {
			It("returns an error", func() {
				err := releaseSource.DownloadReleases("/non-existent-folder", matchedS3Objects, 0)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("/non-existent-folder"))
			})
		})

		Context("when a file can't be downloaded", func() {
			BeforeEach(func() {
				fakeS3Downloader.DownloadCalls(func(w io.WriterAt, i *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
					return 0, errors.New("503 Service Unavailable")
				})
			})

			It("returns an error", func() {
				err := releaseSource.DownloadReleases(releaseDir, matchedS3Objects, 0)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("failed to download file, 503 Service Unavailable\n"))
			})
		})
	})
})
