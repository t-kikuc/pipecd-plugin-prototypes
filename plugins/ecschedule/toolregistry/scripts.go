package toolregistry

const installScript = `
cd {{ .TmpDir }}
{{ if eq .Os "linux" }}
curl -L https://github.com/Songmu/ecschedule/releases/download/{{ .Version }}/ecschedule_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz -o ecschedule_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz
tar -zxvf ecschedule_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz
{{ else }}
echo [DEBUG] https://github.com/Songmu/ecschedule/releases/download/{{ .Version }}/ecschedule_{{ .Version }}_{{ .Os }}_{{ .Arch }}.zip
curl -L https://github.com/Songmu/ecschedule/releases/download/{{ .Version }}/ecschedule_{{ .Version }}_{{ .Os }}_{{ .Arch }}.zip -o ecschedule_{{ .Version }}_{{ .Os }}_{{ .Arch }}.zip
unzip ecschedule_{{ .Version }}_{{ .Os }}_{{ .Arch }}.zip
cd ecschedule_{{ .Version }}_{{ .Os }}_{{ .Arch }}
{{ end }}
mv ecschedule {{ .OutPath }}
`
