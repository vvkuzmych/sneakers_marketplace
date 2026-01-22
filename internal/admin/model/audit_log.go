package model

import (
	"encoding/json"
	"time"
)

// AuditLog represents an admin action audit log entry
type AuditLog struct {
	CreatedAt  time.Time
	Details    map[string]interface{}
	AdminEmail string
	ActionType string
	EntityType string
	IPAddress  string
	ID         int64
	AdminID    int64
	EntityID   int64
}

// AuditAction types
const (
	ActionUserBanned      = "user_banned"
	ActionUserUnbanned    = "user_unbanned"
	ActionUserDeleted     = "user_deleted"
	ActionUserRoleUpdated = "user_role_updated"
	ActionOrderCancelled  = "order_canceled"
	ActionProductFeatured = "product_featured"
	ActionProductHidden   = "product_hidden"
)

// EntityType types
const (
	EntityUser    = "user"
	EntityOrder   = "order"
	EntityProduct = "product"
)

// CreateAuditLogParams contains parameters for creating an audit log
type CreateAuditLogParams struct {
	Details    map[string]interface{}
	ActionType string
	EntityType string
	IPAddress  string
	AdminID    int64
	EntityID   int64
}

// ListAuditLogsParams contains parameters for listing audit logs
type ListAuditLogsParams struct {
	DateFrom   *time.Time
	DateTo     *time.Time
	ActionType string
	AdminID    int64
	Page       int32
	PageSize   int32
}

// DetailsJSON returns details as JSON string
func (a *AuditLog) DetailsJSON() (string, error) {
	if a.Details == nil {
		return "{}", nil
	}
	bytes, err := json.Marshal(a.Details)
	if err != nil {
		return "{}", err
	}
	return string(bytes), nil
}

// ParseDetailsFromJSON parses details from JSON string
func ParseDetailsFromJSON(jsonStr string) (map[string]interface{}, error) {
	var details map[string]interface{}
	if jsonStr == "" || jsonStr == "{}" {
		return make(map[string]interface{}), nil
	}
	err := json.Unmarshal([]byte(jsonStr), &details)
	return details, err
}

// IsValidActionType checks if action type is valid
func IsValidActionType(actionType string) bool {
	validTypes := map[string]bool{
		ActionUserBanned:      true,
		ActionUserUnbanned:    true,
		ActionUserDeleted:     true,
		ActionUserRoleUpdated: true,
		ActionOrderCancelled:  true,
		ActionProductFeatured: true,
		ActionProductHidden:   true,
	}
	return validTypes[actionType]
}

// IsValidEntityType checks if entity type is valid
func IsValidEntityType(entityType string) bool {
	return entityType == EntityUser || entityType == EntityOrder || entityType == EntityProduct
}
