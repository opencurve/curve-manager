/*
*  Copyright (c) 2023 NetEase Inc.
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
 */

/*
* Project: Curve-Manager
* Created Date: 2023-02-11
* Author: wanghai (SeanHai)
 */

package curvebs

import (
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/api/curvebs/manager"
	"github.com/opencurve/curve-manager/api/curvebs/user"
	"github.com/opencurve/pigeon"
)

func NewServer() *pigeon.HTTPServer {
	server := pigeon.NewHTTPServer("curvebs")
	server.Initer(func(cfg *pigeon.Configure) error {
		return core.Init(cfg, server.Logger())
	})
	server.Route("/curvebs",
		core.Rewrite,
		manager.Entrypoint,
		user.Entrypoint)
	server.DefaultRoute(core.Default)
	return server
}

func main() {
	pigeon.Serve(NewServer())
}
