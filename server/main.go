// Copyright (C) 2021 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

//

package main

import (
	"fileservice/pakkusys/application"
	"flag"
)

func main() {
	name := flag.String("name", "fileservice", "App id, default fileservice")
	logger := flag.String("logger", "console", "logger: console, file or unset, default console")
	loglevel := flag.String("loglevel", "debug", "loglevel: debug, info, error or none, default debug")
	logdir := flag.String("logdir", "./logs", "default ./logs/{application name}.log")
	flag.Parse()
	// 启动应用&启用web服务
	application.ApplicationBoot(*name).SetLogger(*logger, *logdir, *loglevel)
	// 启动&&web服务
	application.BootStartWeb(*loglevel == "debug")
}
