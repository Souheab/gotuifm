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
)

func GetFileIcon(fileExtension string) rune {
	switch fileExtension {
	case ".png", ".jpeg":
		return ImageIcon

	case ".mkv", ".mp4", ".webm":
		return VideoIcon

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
