package migrations

import (
	"gin/app/model"

	"gorm.io/gorm"
)

// CreateAgentMessageTable AI消息详情表迁移
type CreateAgentMessageTable struct{}

func (m *CreateAgentMessageTable) ID() string {
	return "20260813_create_agent_message_table"
}

func (m *CreateAgentMessageTable) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.AgentMessage{})
}

func (m *CreateAgentMessageTable) Rollback(db *gorm.DB) error {
	return db.Migrator().DropTable(&model.AgentMessage{})
}
