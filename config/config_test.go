package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestInitConfig(t *testing.T) {
	workDir, _ := os.Getwd()
	path := workDir

	t.Run("正向测试", func(t *testing.T) {
		cfg, err := InitConfig("application", "yaml", path)
		if err != nil {
			t.Fatalf("InitConfig() returned an error: %v, expected no error", err)
		}

		// 断言每一个字段的值是否与配置文件中的一致
		assert.Equal(t, "localhost", cfg.Server.Host)
		assert.Equal(t, "8081", cfg.Server.Port)
		assert.Equal(t, 5*time.Second, cfg.Client.Timeout)
		assert.Equal(t, 3, cfg.Client.MaxRetries)
		assert.Equal(t, 2*time.Second, cfg.Client.Sleeper)
		assert.Equal(t, "https://goerli.blockpi.network/v1/rpc/public", cfg.Client.URL)

	})

	t.Run("反向测试", func(t *testing.T) {
		// 修改配置文件中的一项值，使其无效
		_, err := InitConfig("application", "json", path)
		if err == nil {
			t.Fatal("InitConfig() returned no error, expected an error")
		}
	})
}
