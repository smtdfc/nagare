package upload

import (
	"os"
	"path"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/pkgs/helpers"
	"github.com/smtdfc/nagare/pkgs/paths"
)

type Manager struct {
	logger      *logger.BaseLogger
	attachments map[string]string
}

func (um *Manager) AddAttachment(p string) (string, error) {
	id := helpers.GenerateUUID()
	dir := path.Join(paths.UploadDir, "contents")
	if err := os.MkdirAll(dir, 0755); err != nil {
		um.logger.Error("Failed to get file from form file", "err", err)
		return "", custom_errors.ErrUploadFailed
	}

	uploadPath := path.Join(dir, id)
	if um.attachments == nil {
		um.attachments = make(map[string]string)
	}
	err := helpers.CopyFile(p, uploadPath)
	if err != nil {
		um.logger.Error("Failed to copy file to upload path", "error", err)
		return "", custom_errors.ErrUploadFailed
	}

	um.attachments[id] = uploadPath
	um.logger.Info("Attachment added", "attachment_id", id, "path", uploadPath)

	return id, nil
}

func (um *Manager) GetAttachmentPath(id string) (string, bool) {
	p, exists := um.attachments[id]
	return p, exists
}

func (um *Manager) RemoveAttachment(id string) {
	delete(um.attachments, id)
	um.logger.Info("Attachment removed", "id", id)
}

// @Injectable
func NewUploadManager(logger *logger.BaseLogger) *Manager {
	return &Manager{
		logger: logger.With("module", "upload-manager"),
	}
}
