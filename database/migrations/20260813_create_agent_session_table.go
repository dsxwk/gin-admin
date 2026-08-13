package migrations

import (
	"gin/app/model"
	"gorm.io/gorm"
)

// CreateAgentSessionTable AI会话主表迁移
type CreateAgentSessionTable struct{}

func (m *CreateAgentSessionTable) ID() string {
	return "20260813_create_agent_session_table"
}

func (m *CreateAgentSessionTable) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.AgentSession{})
}

func (m *CreateAgentSessionTable) Rollback(db *gorm.DB) error {
	return db.Migrator().DropTable(&model.AgentSession{})
}
