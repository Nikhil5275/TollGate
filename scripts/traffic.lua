-- Spread benchmark traffic across clients to exercise accepted requests.
local counter = 0
request = function()
  counter = counter + 1
  return wrk.format("GET", "/", {["X-Client-ID"] = "bench-" .. counter})
end
