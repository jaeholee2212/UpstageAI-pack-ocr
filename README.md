# Ocr PoC 
Ocr PoC안의 애플리케이션 서비스들을 위해 인프라 서비스들을 제공합니다. 

![image](https://user-images.githubusercontent.com/90643143/137259358-fccdcb3c-091e-409a-9ca0-57261772976f.png)


제공되는 서비스들은 아래와 같습니다.
1. Load balancer
2. 중앙 로깅 시스템 (A centralized logging system)
3. 데이터베이스들 (Databases)
4. Docker Swarm 관리 툴
5. Feature flagging 툴 (a.k.a an A/B testing)
6. 저장소 서비스들 (Storage services)
7. 모니터링 & 알람 서비스들

각 서비스들은 스택(Stack)이라는 단위로 그룹지어져 있습니다. 개별 스택은 yaml 파일로 작성되었습니다. `yaml`의 내용은 `docker-compose` 파일과 동일하고 추가 `deploy` 섹션이 있습니다.

| 스택 이름 | 용도                                    | yaml                  |
| --------- | --------------------------------------- | --------------------- |
| admin     | Docker swarm 관리 서비스들              | `infra_admin.yaml`    |
| db        | 데이터베이스들                          | `infra_db.yaml`       |
| flags     | 기능 플래그 서비스들. A/B 테스트        | `infra_flags.yaml`    |
| lb        | 로드 발란서 a.k.a reverse-proxy         | `infra_lb.yaml`       |
| logs      | 로깅 시스템                             | `infra_logs.yaml`     |
| mons      | 모니터링 서비스들. 메트릭 서비스들 제공 | `infra_mons.yaml`     |
| registry  | Docker image registry                   | `infra_registry.yaml` |
| storage   | a Blob storage(minio) 제공              | `infra_storage.yaml`  |

# 사용방법
인프라 서비스들을 위해 유틸리티 스크립트를 이용합니다. 스크립트를 실행시킨 컴퓨터에 클러스터를 구성하고 서비스들을 컨테이너 단위로 실행시킵니다.

```
# 시작
./app up

# 종료
./app dn

# 개별 INFRA 스택 시작
./app infra <STACK SHORT NAME> up
# 예
./app infra load-balancer up

# 개별 스택 종료
./app infra load-balancer down

# 개별 스택 다시시작
./app infra load-balancer re
```

# 서비스 접속하기
컴퓨터에서 서비스들의 엔드포인트들을 이해할 수 있게 설정을 바꾸어여합니다.
`/etc/hosts`에 아래 항목들을 추가한 후에 저장합니다.
> 추후에 도메인 설정을 하면 subnames들에 대해서 같은 `ip`를 지정하도록 설정해야합니다.

```
127.0.0.1       traefik.example.com
127.0.0.1       kibana.example.com
127.0.0.1       admin.example.com
127.0.0.1       prometheus.example.com
127.0.0.1       grafana.example.com
127.0.0.1       unleash.example.com
127.0.0.1       portainer.example.com
127.0.0.1       minio.example.com
127.0.0.1       snorkel.example.com
```

브라우져에서 각 서비스들을 접속합니다. 서비스에 대한 항목은 아래와 같습니다.
| 주소                                             | 설명                                               |
| ------------------------------------------------ | -------------------------------------------------- |
| [traefik.example.com](https://traefik.example.com)       | 로드 발란서 관리 툴                                |
| [kibana.example.com](https://kibana.example.com)         | 로그 검색 및 대쉬보드                              |
| [portainer.example.com](https://portainer.example.com)   | Docker Swarm 관리 툴                               |
| [promethues.example.com](https://prometheus.example.com) | 메트릭 탐색 툴                                     |
| [grafana.example.com](https://grafana.example.com)       | 메트릭들 대쉬보드. Prometheus가 메인 데이터 소스임 |
| [unleash.example.com](https://unleash.example.com)       | Feature Flagging 및 A/B 테스트 관리 툴             |
| [minio.example.com](https://minio.example.com)           | Blob 저장소                                        |
| [snorkel.example.com](https://snorkel.example.com)       | Ad-hoc(임시) 데이터 분석 툴                        |

사용자/암호는 모두 `admin`/`12345678`로 통일 되어있습니다. 어떤 서비스는 `admin` 계정 설정을 요구하는데 편하시게 사용하시면 됩니다.

# 환경변수들
환경 변수들은 여러 파일들로 나뉘어져 있습니다. `infra` 스크립트가 취합해서
시스템에 로드합니다.
| 이름         | 설명                                               |
| ------------ | -------------------------------------------------- |
| `.env`       | 개인별 설정을 합니다. 코드에 **체크인 되지 않습니다**. |
| `.env.base`  | 공통된 설정. 체크인 됨.                            |
| `.env.amd64` | x86 계열에 필요한 환경 변수들. 체크인 됨           |
| `.env.arm64` | Apple silicon에 필요한 환경 변수들. 체크인 됨      |

## `.env` 샘플 
```
INFRA_DOMAIN=example.com
INFRA_SWARM_ADVERTISE_ADDR=127.0.0.1
```

# 데모
AWS EC2 머신 한대에 데모 서비스들을 실행시키고 있습니다. 접속 주소는 아래와 같습니다.
> 주의: `/etc/hosts`에 아래 주소들을 반드시 입력해주세요.
```
3.22.209.210    traefik.ec2-3-22-209-210.us-east-2.compute.amazonaws.com
3.22.209.210    kibana.ec2-3-22-209-210.us-east-2.compute.amazonaws.com
3.22.209.210    portainer.ec2-3-22-209-210.us-east-2.compute.amazonaws.com
3.22.209.210    prometheus.ec2-3-22-209-210.us-east-2.compute.amazonaws.com
3.22.209.210    grafana.ec2-3-22-209-210.us-east-2.compute.amazonaws.com
3.22.209.210    unleash.ec2-3-22-209-210.us-east-2.compute.amazonaws.com
3.22.209.210    minio.ec2-3-22-209-210.us-east-2.compute.amazonaws.com
3.22.209.210    snorkel.ec2-3-22-209-210.us-east-2.compute.amazonaws.com
```

접속가능한 데모 서비스들입니다.
| 이름 | 주소 |
|--- |---|
| Load balancer | [https://traefik.ec2-3-22-209-210.us-east-2.compute.amazonaws.com](https://traefik.ec2-3-22-209-210.us-east-2.compute.amazonaws.com) |
| Log viewer | [https://kibana.ec2-3-22-209-210.us-east-2.compute.amazonaws.com](https://kibana.ec2-3-22-209-210.us-east-2.compute.amazonaws.com) |
| Prometheus | [https://prometheus.ec2-3-22-209-210.us-east-2.compute.amazonaws.com](https://prometheus.ec2-3-22-209-210.us-east-2.compute.amazonaws.com) |
| Portainer | [https://portainer.ec2-3-22-209-210.us-east-2.compute.amazonaws.com](https://portainer.ec2-3-22-209-210.us-east-2.compute.amazonaws.com) |
| Grafana | [https://grafana.ec2-3-22-209-210.us-east-2.compute.amazonaws.com](https://grafana.ec2-3-22-209-210.us-east-2.compute.amazonaws.com) |
| Unleash (Feature flagging) | [https://unleash.ec2-3-22-209-210.us-east-2.compute.amazonaws.com](https://unleash.ec2-3-22-209-210.us-east-2.compute.amazonaws.com) |
| Blob storage | [https://minio.ec2-3-22-209-210.us-east-2.compute.amazonaws.com](https://minio.ec2-3-22-209-210.us-east-2.compute.amazonaws.com) |
| Ad-hoc data analysis | [https://snorkel.ec2-3-22-209-210.us-east-2.compute.amazonaws.com](https://snorkel.ec2-3-22-209-210.us-east-2.compute.amazonaws.com) |입

사이트를 클릭하면 경고문구가 나옵니다. SSL 인증서가 Let's encrypt로 되어있는데 인증기관이 제대로 설정이 안되었나봐요. 일단 진행하시면 됩니다.
<img src="https://user-images.githubusercontent.com/90643143/137250193-1c4827ab-ac9d-4c67-9e6c-74d8096fd3a0.png" width="320" >


# Troubleshootings

## Q: `/etc/hosts`를 수정했는데 엔드포인트에 접속하지 못합니다.
원인은 Docker for Mac의 경우 Mac과 Docker를 실행중인 VM에서 네트워크 오류로 발생합니다. 이 경우 Docker에서
제공하는 proxy서버를 이용해서 Mac에서 해당 VM의 특정 컨테이너에 접속하게 합니다.

### 스텝1 - Docker에게 proxy 관련 설정을 합니다
```
cd ~/Library/Group\ Containers/group.com.docker/
mv settings.json settings.json.backup
cat settings.json.backup | jq '.["socksProxyPort"]=8888' > settings.json
```
### 스텝2 - Docker for Mac을 재시작합니다
이제 Docker가 proxy 서버를 실행시킵니다.

### 스텝3 - Mac에서 네트워크 proxy를 지정합니다
`System Preferences` -> `Network` -> `Advanced` -> `Proxies` 탭 에서 아래 그림과 같이 설정합니다.
![F1D5C125-A26C-4729-B92A-540870045297](https://user-images.githubusercontent.com/90643143/137250739-89e1b0f8-d8f8-4aac-8973-69bc05713e77.png)

한가지 단점은, proxy를 위해서 Docker for Mac이 항상 실행 중이어야 합니다.

소스 - https://github.com/docker/for-mac/issues/2670#issuecomment-372365274





