# Yungang Grottoes Condensation Monitoring

云冈石窟第 9、10 窟凝结水监测与治理协同服务。系统接收窟内外温湿度与露点观测，维护监测点和治理设备批次，协调治理计划、设备调度、异常事件处置和审计留痕。响应计划、容量约束、请求分配和调度签收是对“看得清、测得准、治得好”流程的后端化实现。

## Run

`GOTOOLCHAIN=local go run ./cmd/server`

The service exposes `/healthz`, `/readyz`, `/v1/login`, `/v1/plans`, `/v1/requests`, `/v1/dispatches`, `/v1/incidents` and audit endpoints. 角色包括 coordinator（治理计划）、dispatcher（设备调度）和 auditor（审计读取）。

SQLite 数据库在启动时从空库执行版本化 migration，并在重启后恢复状态。coordinator 可以激活治理计划并批准设备请求；dispatcher 推进设备投放和签收；auditor 读取不可变审计事件。所有跨实体写入使用事务、幂等键和版本控制，后台 worker 负责重试和永久失败记录。

## Verify

`go test ./... -count=1`

`go test -race ./... -count=1`

`go vet ./...`

`go build ./...`
