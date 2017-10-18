package commands

import "github.com/pivotal-cf/jhanda/flags"

type BakeConfig struct {
	ReleaseDirectories   flags.StringSlice `short:"rd"   long:"releases-directory"      description:"path to the release tarballs directory"`
	MigrationDirectories flags.StringSlice `short:"md"   long:"migrations-directory"    description:"path to the migrations directory"`
	EmbedPaths           flags.StringSlice `short:"e"    long:"embed"                   description:"path to additional files to embed"`
	StemcellTarball      string            `short:"st"   long:"stemcell-tarball"        description:"location of the stemcell tarball"`
	Metadata             string            `short:"m"    long:"metadata"                description:"location of the metadata file"`
	Version              string            `short:"v"    long:"version"                 description:"version of the tile"`
	ProductName          string            `short:"pn"   long:"product-name"            description:"product name"`
	OutputFile           string            `short:"o"    long:"output-file"             description:"the output path of the tile"`
	StubReleases         bool              `short:"sr"   long:"stub-releases"           description:"don't include release tarballs"`
}
