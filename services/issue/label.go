// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issue

import (
	"code.gitea.io/gitea/models"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/notification"
)

// ClearLabels clears all of an issue's labels
func ClearLabels(issue *models.Issue, doer *user_model.User) (err error) {
	if err = models.ClearIssueLabels(issue, doer); err != nil {
		return
	}

	notification.NotifyIssueClearLabels(doer, issue)

	return nil
}

// AddLabel adds a new label to the issue.
func AddLabel(issue *models.Issue, doer *user_model.User, label *models.Label) error {
	if err := models.NewIssueLabel(issue, label, doer); err != nil {
		return err
	}

	notification.NotifyIssueChangeLabels(doer, issue, []*models.Label{label}, nil)
	return nil
}

// AddLabels adds a list of new labels to the issue.
func AddLabels(issue *models.Issue, doer *user_model.User, labels []*models.Label) error {
	if err := models.NewIssueLabels(issue, labels, doer); err != nil {
		return err
	}

	notification.NotifyIssueChangeLabels(doer, issue, labels, nil)
	return nil
}

// RemoveLabel removes a label from issue by given ID.
func RemoveLabel(issue *models.Issue, doer *user_model.User, label *models.Label) error {
	if err := issue.LoadRepo(); err != nil {
		return err
	}

	perm, err := models.GetUserRepoPermission(issue.Repo, doer)
	if err != nil {
		return err
	}
	if !perm.CanWriteIssuesOrPulls(issue.IsPull) {
		if label.OrgID > 0 {
			return models.ErrOrgLabelNotExist{}
		}
		return models.ErrRepoLabelNotExist{}
	}

	if err := models.DeleteIssueLabel(issue, label, doer); err != nil {
		return err
	}

	notification.NotifyIssueChangeLabels(doer, issue, nil, []*models.Label{label})
	return nil
}

// ReplaceLabels removes all current labels and add new labels to the issue.
func ReplaceLabels(issue *models.Issue, doer *user_model.User, labels []*models.Label) error {
	old, err := models.GetLabelsByIssueID(issue.ID)
	if err != nil {
		return err
	}

	if err := models.ReplaceIssueLabels(issue, labels, doer); err != nil {
		return err
	}

	notification.NotifyIssueChangeLabels(doer, issue, labels, old)
	return nil
}
