package toolregistry

const installScript = `
cd {{ .TmpDir }}
curl -L https://github.com/kayac/ecspresso/releases/download/v{{ .Version }}/ecspresso_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz -o ecspresso_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz
tar -zxvf ecspresso_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz
mv ecspresso {{ .OutPath }}
`
