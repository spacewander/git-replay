-- luacheck: globals git_replay
local commit = git_replay.get_commit()
local dump = "fixture/dump_commit.actual.txt"
local f = io.open(dump, 'w')

local output = {}
for attr, value in pairs(commit) do
    if attr ~= 'tags' then
        table.insert(output, attr..': '..value)
    end
end

for _, tag in ipairs(commit.tags) do
    for attr, value in pairs(tag) do
        table.insert(output, 'tag: '..attr..': '..value)
    end
end
table.sort(output)

f:write(table.concat(output, '\n'))
f:close()
