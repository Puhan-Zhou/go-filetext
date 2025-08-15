# Go plaintext getter
A library to get plaintext from different file type.

## Support file type
- pdf
- docx (modern Word format)
- xlsx (modern Excel format)
- pptx (modern PowerPoint format)
- plain text file(like txt, yaml, csv, markdown, ...)
- (future support) jpeg, png or other picture format, use ocr

### Legacy Office Format Support
Legacy Microsoft Office formats (DOC, XLS, PPT) are detected and recognized by the library, but text extraction is not supported due to their proprietary binary format. When these formats are encountered, the library will provide a helpful error message suggesting conversion to their modern equivalents (DOCX, XLSX, PPTX).