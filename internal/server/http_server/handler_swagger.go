package http_server

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

var swaggerHTML = `
<html>
    <head>
        <script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
        <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3/swagger-ui.css"/>
        <title>ssup2ket auth service</title>
    </head>
    <body>
        <div id="swagger-ui"></div> 
        <script>
            window.onload = function () {
                const ui = SwaggerUIBundle({
                    url: "{{ .URL }}/v1/swagger/spec",
                    dom_id: '#swagger-ui',
                    deepLinking: true,
                    presets: [
                        SwaggerUIBundle.presets.apis,
                        SwaggerUIBundle.SwaggerUIStandalonePreset
                    ],
                    plugins: [
                        SwaggerUIBundle.plugins.DownloadUrl
                    ],
                })
                window.ui = ui
            }
        </script>
    </body>
</html>
`

func getSwaggerSpecHandler(url string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get swagger spec
		swagger, err := GetSwagger()
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to get swagger spec")
			_ = render.Render(w, r, getErrRendererServerError())
			return
		}

		// Validate swagger spec
		if len(swagger.Servers) != 1 {
			log.Ctx(ctx).Error().Msg("No server info in the swagger spec")
			_ = render.Render(w, r, getErrRendererServerError())
			return
		}

		// Set server URL
		swagger.Servers[0].URL = url
		render.JSON(w, r, swagger)
	}

	return fn
}

func getSwaggerUIHandler(url string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Set server URL
		_, err := w.Write([]byte(strings.ReplaceAll(swaggerHTML, "{{ .URL }}", url)))
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to get swagger UI")
			_ = render.Render(w, r, getErrRendererServerError())
			return
		}
	}
	return fn
}
