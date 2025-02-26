package toolregistry

const installScript = `
cd {{ .TmpDir }}
curl -L https://github.com/fujiwara/lambroll/releases/download/{{ .Version }}/lambroll_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz -o lambroll_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz
tar -zxvf lambroll_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz
mv lambroll {{ .OutPath }}
`
