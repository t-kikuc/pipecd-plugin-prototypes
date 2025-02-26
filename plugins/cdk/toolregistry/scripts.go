package toolregistry

const installScript = `
cd {{ .TmpDir }}
curl -L https://github.com/fujiwara/cdk/releases/download/{{ .Version }}/cdk_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz -o cdk_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz
tar -zxvf cdk_{{ .Version }}_{{ .Os }}_{{ .Arch }}.tar.gz
mv cdk {{ .OutPath }}
`
