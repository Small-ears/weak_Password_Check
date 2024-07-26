1、检测出所有的弱口令，空口令

2、程序中存在一处bug（已解决），当存在一些mysql账户权限不足时则无法执行db.Ping(),导致出现漏报，通过增加检查err中是否包含strings.Contains(err.Error(), "to database")，对存在的结果输出到保存到文件中，但是Result为flase，建议扫描后使用手工对结果文件中为false的结果进行验证。