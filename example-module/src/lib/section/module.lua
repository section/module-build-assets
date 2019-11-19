local envvar = require "section.environment_variables"

local setmetatable = setmetatable
local tonumber = tonumber
local tostring = tostring
local rawget = rawget

local _M = {
    _VERSION = "1.0.0"
}

local mt = { __index = _M }

-- Method constants
local methods = {
 GET = ngx.HTTP_GET,
 HEAD = ngx.HTTP_HEAD,
 PUT = ngx.HTTP_PUT,
 POST = ngx.HTTP_POST,
 DELETE = ngx.HTTP_DELETE,
 OPTIONS = ngx.HTTP_OPTIONS,
 MKCOL = ngx.HTTP_MKCOL,
 COPY  = ngx.HTTP_COPY,
 MOVE = ngx.HTTP_MOVE,
 PROPFIND = ngx.HTTP_PROPFIND,
 PROPATC = ngx.HTTP_PROPPATC,
 LOCK = ngx.HTTP_LOCK,
 UNLOCK = ngx.HTTP_UNLOCK,
 PATCH = ngx.HTTP_PATCH,
 TRACE = ngx.HTTP_TRACE
}

-- helper function to concatenate tables
local function TableConcat(t1, t2)
  for i=1,#t2 do
      t1[#t1+1] = t2[i]
  end
  return t1
end

-- Set response headers for response being send downstream
local function set_resp_headers(resp_headers)
  -- Uncomment if local set cookie is used
  -- temp_local_setcookie = ngx.header["Set-Cookie"]
  setcookie = {}
  for k, v in pairs(resp_headers) do
    -- Set-Cookie is a snowflake - when doing sub-requests simply copying the header value can cause issues if there are
    -- variations in the capitalisation. This should NOT be the case as the RFC says headers are case-insensitive.
    -- However, we have seen this behaviour so we capture all the variations as lower-case and add them to a table.
    -- A check is also done on the returned header value to see if it is a table and unpack it if it is.

    canonicalKey = string.lower(k)
    if (canonicalKey == "set-cookie") then
      if type(v) == "table" then
        setcookie = TableConcat(setcookie, v)
      else
        table.insert(setcookie, v)
      end
    else
      ngx.header[k] = v
    end
  end
  -- Uncomment if local set cookie is used
  -- if type(temp_local_setcookie) == "table" then
  --  setcookie = TableConcat(setcookie, temp_local_setcookie)
  -- else
  --   table.insert(setcookie, temp_local_setcookie)
  -- end
  ngx.header["Set-Cookie"] = setcookie
end


-- Pass the request upstream
-- Similar logic can be used to pass the request to an API endpoint.
local function pass_request()
  ngx.req.read_body()
  local resp, err = ngx.location.capture("/.well-known/section-io/examplemodule/pass-request" .. ngx.var.request_uri, { method = methods[ngx.var.request_method]})

  set_resp_headers(resp.header)
  return resp.status, resp
end

-- Make decision for the incoming request
function _M.custom_lua(self, url)
  -- Pass the request to upstream if the module is not enabled.
  if envvar.ENABLED == false then
    local status, response = pass_request()
    return status, response.body
  end


  if not ngx.req.get_headers()["X-Forwarded-Proto"] then
    -- Assume HTTP if no X-Forwarded-Proto header
    ngx.req.set_header("X-Forwarded-Proto", "http")
  end

  if not ngx.req.get_headers()["Host"] then
    return 400, "Bad request: No Host header"
  end

  -- Fetch the response body from upstream and deliver
  local status, response = pass_request()
  -- Custom logic
  ngx.header["x-module-header"] = "Section module";

  return status, response.body
end

return _M
