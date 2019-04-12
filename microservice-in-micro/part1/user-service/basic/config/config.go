package config

import (
	"fmt"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	err error
)

var (
	defaultRootPath         = "app"
	defaultConfigFilePrefix = "application-"
	consulConfig            defaultConsulConfig
	mysqlConfig             defaultMysqlConfig
	profiles                defaultProfiles
	m                       sync.RWMutex
	inited                  bool
)

// InitConfig 初始化配置
func InitConfig() {

	m.Lock()
	defer m.Unlock()

	if inited {
		log.Fatal(fmt.Errorf("[InitConfig] 配置已经初始化过"))
		return
	}

	// 加载yml配置
	// 先加载基础配置
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("./", string(filepath.Separator))))

	pt := filepath.Join(appPath, "conf")
	os.Chdir(appPath)

	// 找到application.yml文件
	if err = config.Load(file.NewSource(file.WithPath(pt + "/application.yml"))); err != nil {
		panic(err)
	}

	// 找到需要引入的新配置文件
	if err = config.Get(defaultRootPath, "profiles").Scan(&profiles); err != nil {
		panic(err)
	}

	log.Logf("[InitConfig] 加载配置文件：path: %s, %+v\n", pt+"/application.yml", profiles)

	// 开始导入新文件
	if len(profiles.GetInclude()) > 0 {
		include := strings.Split(profiles.GetInclude(), ",")

		sources := make([]source.Source, len(include))
		for i := 0; i < len(include); i++ {
			filePath := pt + string(filepath.Separator) + defaultConfigFilePrefix + strings.TrimSpace(include[i]) + ".yml"
			fmt.Printf(filePath + "\n")
			sources[i] = file.NewSource(file.WithPath(filePath))
		}

		// 加载include的文件
		if err = config.Load(sources...); err != nil {
			panic(err)
		}
	}

	// 赋值
	config.Get(defaultRootPath, "consul").Scan(&consulConfig)
	config.Get(defaultRootPath, "mysql").Scan(&mysqlConfig)

	// 标记已经初始化
	inited = true
}

// GetMysqlConfig 获取mysql配置
func GetMysqlConfig() (ret MysqlConfig) {
	return mysqlConfig
}

// GetConsulConfig 获取Consul配置
func GetConsulConfig() (ret ConsulConfig) {
	return consulConfig
}
