package toolregistry

const nodeInstallScript = `
cd {{ .TmpDir }}
curl -L https://nodejs.org/dist/{{ .Version }}/node-{{ .Version }}-{{ .Os }}-{{ .Arch }}.tar.gz -o node-{{ .Version }}-{{ .Os }}-{{ .Arch }}.tar.gz
tar -zxf node-{{ .Version }}-{{ .Os }}-{{ .Arch }}.tar.gz
cp -L node-{{ .Version }}-{{ .Os }}-{{ .Arch }}/bin/npm {{ .OutPath }}
`

// const nodeInstallScript = `
// cd {{ .TmpDir }}
// curl -L https://nodejs.org/dist/{{ .Version }}/node-{{ .Version }}-{{ .Os }}-{{ .Arch }}.tar.gz -o {{ .OutPath }}
// `

// mv node-{{ .Version }}-{{ .Os }}-{{ .Arch }}/bin/node {{ .OutPath }}

// const nodeInstallScript = `
// echo "hoge" > {{ .OutPath }}
// `

// const cdkInstallScript = `
// cd {{ .TmpDir }}
// npm install -g aws-cdk@{{ .Version }}
// cp -L $(npm root -g)/aws-cdk/bin/cdk {{ .OutPath }}
// `
