package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
)

const (
	vclCompileCommandline = `"{{.bsdDir}}\bin\dcc{{.arch}}.exe" -$Z4 {{if eq .arch "64"}}-$B+ {{end}}-$D0 -$L- -$Y- --no-config -B -Q -TX.dll -AGenerics.Collections=System.Generics.Collections;Generics.Defaults=System.Generics.Defaults;WinTypes=Winapi.Windows;WinProcs=Winapi.Windows;DbiTypes=BDE;DbiProcs=BDE;DbiErrs=BDE -DLIBVCL;RELEASE -E"{{.binFileDir}}" -I"{{.bsdDir}}\lib\Win{{.arch}}\release";{{.UserDir}}\Imports;"{{.bsdDir}}\Imports";{{.UserDir}}\Dcp;"{{.bsdDir}}\include"; -LE{{.UserDir}}\Bpl -LN{{.UserDir}}\Dcp -NU"{{.objFileDir}}\Win{{.arch}}\Release" -NSWinapi;System.Win;Data.Win;Datasnap.Win;Web.Win;Soap.Win;Xml.Win;Bde;System;Xml;Data;Datasnap;Web;Soap;Vcl;Vcl.Imaging;Vcl.Touch;Vcl.Samples;Vcl.Shell; -O"{{.bsdDir}}\lib\Win{{.arch}}\release";{{.UserDir}}\Imports;"{{.bsdDir}}\Imports";{{.UserDir}}\Dcp;"{{.bsdDir}}\include"; -R"{{.bsdDir}}\lib\Win{{.arch}}\release";{{.UserDir}}\Imports;"{{.bsdDir}}\Imports";{{.UserDir}}\Dcp;"{{.bsdDir}}\include"; -U"{{.bsdDir}}\lib\Win{{.arch}}\release";{{.UserDir}}\Imports;"{{.bsdDir}}\Imports";{{.UserDir}}\Dcp;"{{.bsdDir}}\include"; -NB{{.UserDir}}\Dcp -NH{{.UserDir}}\hpp\Win{{.arch}} -NO"{{.objFileDir}}\Win{{.arch}}\Release" "{{.projectFileName}}"&&Move /Y "{{.binFileDir}}\vcl.dll" "{{.binFileDir}}\libvcl{{if eq .arch "64"}}x64{{end}}.dll"`
)

func buildVCL(bsdDir, UserDir, arch, projectFileName, objFileDir, binFileDir string) error {

	// lcl.lpi 这里面要修改版本号之类的
	// ProductVersion="1.2.8.0"

	if !checkFileExists(bsdDir) {
		return errors.New("未找到bsd目录。")
	}
	if !checkFileExists(projectFileName) {
		return errors.New("libvcl工程文件未找到。")
	}

	bsdDir = fixDirName(bsdDir)
	objFileDir = fixDirName(objFileDir)
	binFileDir = fixDirName(binFileDir)
	UserDir = fixDirName(UserDir)

	fmt.Println("bsdDir:", bsdDir)
	fmt.Println("projectFileName: ", projectFileName)
	fmt.Println("objFileDir: ", objFileDir)
	fmt.Println("binFileDir: ", binFileDir)

	var buf bytes.Buffer
	tmp := template.New("lcl")
	tmp.Parse(vclCompileCommandline)
	tmp.Execute(&buf, map[string]interface{}{
		"bsdDir":          bsdDir,
		"arch":            arch, // 32 64
		"projectFileName": projectFileName,
		"objFileDir":      objFileDir, // 编译时中间文件目录
		"binFileDir":      binFileDir,
		"UserDir":         UserDir,
	})

	// libvcl二进制输出目录不存在，则创建
	if !checkFileExists(binFileDir) {
		if err := os.MkdirAll(binFileDir, 0775); err != nil {
			return err
		}
	}
	//
	//// 中间文件输出
	libOutDir := objFileDir + "\\Win" + arch + "\\Release"
	if !checkFileExists(libOutDir) {
		if err := os.MkdirAll(libOutDir, 0775); err != nil {
			return err
		}
	}

	fmt.Println(buf.String())
	return nil
}