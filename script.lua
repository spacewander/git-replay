-- luacheck: globals git_replay
local commit = git_replay.get_commit()
local story = {}
story[1] = ("%s committed %s in %s with title %s"):format(
    commit.name, commit.hash:sub(1, 7), commit.date, commit.title)
if commit.tags then
    local tag_names = {}
    for i = 1, #commit.tags do
        table.insert(tag_names, commit.tags[i].name)
    end
    story[1] = story[1] .. "\ntags: " .. table.concat(tag_names, ' ')
end

-- date format: 2017-03-22 16:53:37 +0800
local cur_month = commit.date:sub(6, 7)
local month = git_replay.month
if not month then
    git_replay.month = {}
    month = git_replay.month
    month.last_month = cur_month
    month[cur_month] = {}
end

if month[cur_month] then
    table.insert(month[cur_month], commit)
else
    local commit_num = #month[month.last_month]
    story[2] = ("last month: %d commits were push to repo in all branches"):format(commit_num)
    month[month.last_month] = nil
    month.last_month = cur_month
    month[cur_month] = {commit}
end

git_replay.display(table.concat(story, '\n'))
