package download

type DownloaderUsecases struct {
}

func NewDownloadUsecases() Downloader {
	return &DownloaderUsecases{}
}
