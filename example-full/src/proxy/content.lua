-- import the module configuration, lib/section/environment_variables.luaass
local envvar = require "section.environment_variables"
-- import the module lua file, lib/section/module.lua
local module = require "section.module"

-- call the function in lib/section/module.lua
local status, results = module:custom_lua(
    ngx.var.request_uri
)

-- Deliver the response downstream
ngx.status = status
ngx.print(results)
ngx.eof()

return