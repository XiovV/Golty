package ytdl

type Ytdl struct {
	BaseCommand string
}

func New(baseCommand string) *Ytdl {
	return &Ytdl{BaseCommand: baseCommand}
}
