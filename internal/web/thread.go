// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/superseriousbusiness/gotosocial/internal/ap"
	apimodel "github.com/superseriousbusiness/gotosocial/internal/api/model"
	apiutil "github.com/superseriousbusiness/gotosocial/internal/api/util"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
)

func (m *Module) threadGETHandler(c *gin.Context) {
	ctx := c.Request.Context()

	authed, err := oauth.Authed(c, false, false, false, false)
	if err != nil {
		apiutil.WebErrorHandler(c, gtserror.NewErrorUnauthorized(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	// usernames on our instance will always be lowercase
	username := strings.ToLower(c.Param(usernameKey))
	if username == "" {
		err := errors.New("no account username specified")
		apiutil.WebErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	// status ids will always be uppercase
	statusID := strings.ToUpper(c.Param(statusIDKey))
	if statusID == "" {
		err := errors.New("no status id specified")
		apiutil.WebErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	instance, err := m.processor.InstanceGetV1(ctx)
	if err != nil {
		apiutil.WebErrorHandler(c, gtserror.NewErrorInternalError(err), m.processor.InstanceGetV1)
		return
	}

	instanceGet := func(ctx context.Context) (*apimodel.InstanceV1, gtserror.WithCode) {
		return instance, nil
	}

	// do this check to make sure the status is actually from a local account,
	// we shouldn't render threads from statuses that don't belong to us!
	if _, errWithCode := m.processor.Account().GetLocalByUsername(ctx, authed.Account, username); errWithCode != nil {
		apiutil.WebErrorHandler(c, errWithCode, instanceGet)
		return
	}

	status, errWithCode := m.processor.Status().Get(ctx, authed.Account, statusID)
	if errWithCode != nil {
		apiutil.WebErrorHandler(c, errWithCode, instanceGet)
		return
	}

	if !strings.EqualFold(username, status.Account.Username) {
		err := gtserror.NewErrorNotFound(errors.New("path username not equal to status author username"))
		apiutil.WebErrorHandler(c, gtserror.NewErrorNotFound(err), instanceGet)
		return
	}

	// if we're getting an AP request on this endpoint we
	// should render the status's AP representation instead
	accept := c.NegotiateFormat(string(apiutil.TextHTML), string(apiutil.AppActivityJSON), string(apiutil.AppActivityLDJSON))
	if accept == string(apiutil.AppActivityJSON) || accept == string(apiutil.AppActivityLDJSON) {
		m.returnAPStatus(ctx, c, username, statusID, accept)
		return
	}

	context, errWithCode := m.processor.Status().ContextGet(ctx, authed.Account, statusID)
	if errWithCode != nil {
		apiutil.WebErrorHandler(c, errWithCode, instanceGet)
		return
	}

	stylesheets := []string{
		assetsPathPrefix + "/Fork-Awesome/css/fork-awesome.min.css",
		distPathPrefix + "/status.css",
	}
	if config.GetAccountsAllowCustomCSS() {
		stylesheets = append(stylesheets, "/@"+username+"/custom.css")
	}

	c.HTML(http.StatusOK, "thread.tmpl", gin.H{
		"instance":    instance,
		"status":      status,
		"context":     context,
		"ogMeta":      ogBase(instance).withStatus(status),
		"stylesheets": stylesheets,
		"javascript":  []string{distPathPrefix + "/frontend.js"},
	})
}

func (m *Module) returnAPStatus(ctx context.Context, c *gin.Context, username string, statusID string, accept string) {
	verifier, signed := c.Get(string(ap.ContextRequestingPublicKeyVerifier))
	if signed {
		ctx = context.WithValue(ctx, ap.ContextRequestingPublicKeyVerifier, verifier)
	}

	signature, signed := c.Get(string(ap.ContextRequestingPublicKeySignature))
	if signed {
		ctx = context.WithValue(ctx, ap.ContextRequestingPublicKeySignature, signature)
	}

	status, errWithCode := m.processor.Fedi().StatusGet(ctx, username, statusID)
	if errWithCode != nil {
		apiutil.WebErrorHandler(c, errWithCode, m.processor.InstanceGetV1) //nolint:contextcheck
		return
	}

	b, mErr := json.Marshal(status)
	if mErr != nil {
		err := fmt.Errorf("could not marshal json: %s", mErr)
		apiutil.WebErrorHandler(c, gtserror.NewErrorInternalError(err), m.processor.InstanceGetV1) //nolint:contextcheck
		return
	}

	c.Data(http.StatusOK, accept, b)
}
