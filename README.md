## Readme

run command on slack:
```
/pr [user] [repo-name] [limit]
```

sample
```
/pr octocat Hello-World 10
```

Running Server

```
go run cmd/api/main.go [User Token]
```

## Changelog
=== initial commit
- baru bisa print array aja di slack
- repo name masih hardcode
- belum ada list PR

=== #2
- Response message lebih rapih
- pake attachment untuk list PR

=== #3
- node list fix
