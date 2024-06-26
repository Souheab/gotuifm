package main

const (
	FolderIcon   = ''
	FileIcon     = ''
	ImageIcon    = ''
	GoLangIcon   = ''
	ConfigIcon   = ''
	OrgIcon      = ''
	MarkdownIcon = ''
	CIcon        = ''
	VideoIcon    = ''
	PDFIcon      = ''
)

func GetFileIcon(fileExtension string) rune {
	if IsImageFileExtension(fileExtension) {
		return ImageIcon
	}

	switch fileExtension {
	case ".mkv", ".mp4", ".webm":
		return VideoIcon

  // Documents
	case ".pdf":
		return PDFIcon

	// Config Files
	case ".mod":
		return ConfigIcon

	// Markup languages
	case ".org":
		return OrgIcon
	case ".md":
		return MarkdownIcon

		// Programming Languages
	case ".go":
		return GoLangIcon
	case ".c":
		return CIcon
	}

	return FileIcon
}
