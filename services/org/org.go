// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package org

import (
	"fmt"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/models/organization"
	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/storage"
	"code.gitea.io/gitea/modules/util"
)

// DeleteOrganization completely and permanently deletes everything of organization.
func DeleteOrganization(org *organization.Organization) error {
	ctx, commiter, err := db.TxContext()
	if err != nil {
		return err
	}
	defer commiter.Close()

	// Check ownership of repository.
	count, err := repo_model.GetRepositoryCount(ctx, org.ID)
	if err != nil {
		return fmt.Errorf("GetRepositoryCount: %v", err)
	} else if count > 0 {
		return models.ErrUserOwnRepos{UID: org.ID}
	}

	if err := organization.DeleteOrganization(ctx, org); err != nil {
		return fmt.Errorf("DeleteOrganization: %v", err)
	}

	if err := commiter.Commit(); err != nil {
		return err
	}

	// FIXME: system notice
	// Note: There are something just cannot be roll back,
	//	so just keep error logs of those operations.
	path := user_model.UserPath(org.Name)

	if err := util.RemoveAll(path); err != nil {
		return fmt.Errorf("Failed to RemoveAll %s: %v", path, err)
	}

	if len(org.Avatar) > 0 {
		avatarPath := org.CustomAvatarRelativePath()
		if err := storage.Avatars.Delete(avatarPath); err != nil {
			return fmt.Errorf("Failed to remove %s: %v", avatarPath, err)
		}
	}

	return nil
}
