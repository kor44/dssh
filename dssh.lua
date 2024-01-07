local parser = clink.argmatcher("dssh")

local function get_hosts()
    local handle = io.popen("dssh --list-hosts")
    local result = handle:read("*all")
    handle:close()

    local hosts = {}
    for host in string.gmatch(result, "%S+") do
      host = string.gsub(host, "  ", "")
      table.insert(hosts, host)
    end
    return hosts
end

parser:addarg({get_hosts})
parser:addflags("--show", "--install-completion", "-h", "--help")
parser:nofiles()