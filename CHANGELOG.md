# Changelog

## Next release after 1.200.0 (2026-11-16): three operations join the section they belong to

Breaking for code calling the three operations below. Versions follow the
API's version.

The API moved three operations into the section, and the service, they are
about. Commands are generated from the sections and operation IDs, so regenerating
(`scripts/generate-commands.sh`) renames two commands and moves one path:

| Before | After |
|---|---|
| `linkbreakers workspace-tokens list-api-scopes` | `linkbreakers workspace-tokens list-api-scopes` (same command, new path) |
| `linkbreakers integrations list-for-link` | `linkbreakers links list-integrations` |
| `linkbreakers auth request-email-verification` | `linkbreakers members request-email-verification` |

| Operation | Path |
|---|---|
| List the scopes a token can have | `GET /v1/workspace-tokens/scopes` (was `GET /v1/api-scopes`, removed with no alias: no integration called it in the last 14 days) |
| List the integrations active on a link | `GET /v1/links/{linkId}/integrations` (unchanged) |
| Resend a member's verification email | `POST /v1/members/email-verification` (unchanged) |

The `integrations` group is gone; the integration groups for Slack, forms,
contact tools, Google Analytics and retargeting pixels stay.

## 1.198.1 to 1.200.0 (2026-11-15 and 2026-11-16): the public API contract cleanup

Breaking for code written against 1.198.0 or earlier. Versions follow the
API's version. The renamed and removed classes, operations and paths below
shipped in 1.198.1; the verification email moved to Members in 1.199.0; the
former paths were removed and Google Forms became one of the form tools in
1.200.0.

The API renamed its services and operations so the contract reads the way it
works. Commands are generated from those names, so regenerating
(`scripts/generate-commands.sh`) renames these commands:

| Before | After |
|---|---|
| `linkbreakers allowances list-allowances` | `linkbreakers plan-usage list` |
| `linkbreakers analytics-dashboard create-dashboard` | `linkbreakers analytics-dashboards create` |
| `linkbreakers analytics-dashboard delete-dashboard` | `linkbreakers analytics-dashboards delete` |
| `linkbreakers analytics-dashboard generate-widget-summary` | `linkbreakers analytics-dashboards generate-widget-summary` |
| `linkbreakers analytics-dashboard get-widget-summary` | `linkbreakers analytics-dashboards get-widget-summary` |
| `linkbreakers analytics-dashboard list-dashboards` | `linkbreakers analytics-dashboards list` |
| `linkbreakers analytics-dashboard update-dashboard` | `linkbreakers analytics-dashboards update` |
| `linkbreakers analytics-data get-attribute-distribution` | `linkbreakers analytics get-attribute-distribution` |
| `linkbreakers analytics-data get-attribute-metric` | `linkbreakers analytics get-attribute-metric` |
| `linkbreakers analytics-data get-attribute-values` | `linkbreakers analytics get-attribute-values` |
| `linkbreakers analytics-data get-breakdown` | `linkbreakers analytics get-breakdown` |
| `linkbreakers analytics-data get-device-distribution` | `linkbreakers analytics get-device-distribution` |
| `linkbreakers analytics-data get-geo-breakdown` | `linkbreakers analytics get-geo-breakdown` |
| `linkbreakers analytics-data get-grouped-time-series` | `linkbreakers analytics get-grouped-time-series` |
| `linkbreakers analytics-data get-hour-of-day-distribution` | `linkbreakers analytics get-hour-of-day-distribution` |
| `linkbreakers analytics-data get-shared-locations` | `linkbreakers analytics get-shared-locations` |
| `linkbreakers analytics-data get-step-funnel` | `linkbreakers analytics get-step-funnel` |
| `linkbreakers analytics-data get-top-links` | `linkbreakers analytics get-top-links` |
| `linkbreakers analytics-data get-traffic-split-results` | `linkbreakers analytics get-traffic-split-results` |
| `linkbreakers analytics-data get-visitor-trends` | `linkbreakers analytics get-visitor-trends` |
| `linkbreakers analytics-data list-visitor-countries` | `linkbreakers analytics list-visitor-countries` |
| `linkbreakers custom-domains check` | `linkbreakers custom-domains verify` |
| `linkbreakers events list-event-traces` | `linkbreakers events list-traces` |
| `linkbreakers events list-events` | `linkbreakers events list` |
| `linkbreakers links shortlink-availability` | `linkbreakers links check-shortlink-availability` |
| `linkbreakers links shortlink-random` | `linkbreakers links generate-shortlink` |
| `linkbreakers metrics get-events-metric` | `linkbreakers analytics get-events-metric` |
| `linkbreakers metrics get-workspace-metric` | `linkbreakers analytics get-workspace-metric` |
| `linkbreakers page-theme-templates create-page-theme-template` | `linkbreakers page-theme-templates create` |
| `linkbreakers page-theme-templates delete-page-theme-template` | `linkbreakers page-theme-templates delete` |
| `linkbreakers page-theme-templates get-page-theme-template` | `linkbreakers page-theme-templates get` |
| `linkbreakers page-theme-templates import-page-theme-template` | `linkbreakers page-theme-templates import` |
| `linkbreakers page-theme-templates list-page-theme-templates` | `linkbreakers page-theme-templates list` |
| `linkbreakers page-theme-templates public-get-page-theme-template` | `linkbreakers page-theme-templates get-public` |
| `linkbreakers page-theme-templates public-list-page-theme-templates` | `linkbreakers page-theme-templates list-public` |
| `linkbreakers page-theme-templates publish-page-theme-template` | `linkbreakers page-theme-templates publish` |
| `linkbreakers page-theme-templates update-page-theme-template` | `linkbreakers page-theme-templates update` |
| `linkbreakers page-theme-templates update-page-theme-template-styling` | `linkbreakers page-theme-templates update-styling` |
| `linkbreakers page-themes apply-page-theme-template` | `linkbreakers page-themes apply-template` |
| `linkbreakers page-themes bulk-update-page-themes` | `linkbreakers page-themes bulk-update` |
| `linkbreakers page-themes get-page-theme` | `linkbreakers page-themes get` |
| `linkbreakers page-themes list-page-themes` | `linkbreakers page-themes list` |
| `linkbreakers page-themes update-page-theme` | `linkbreakers page-themes update` |
| `linkbreakers pricing list-pricing-plans` | `linkbreakers pricing list-plans` |
| `linkbreakers qrcode-design apply-template` | `linkbreakers qrcode-designs apply-template` |
| `linkbreakers qrcode-design delete` | `linkbreakers qrcode-designs delete` |
| `linkbreakers qrcode-design get` | `linkbreakers qrcode-designs get` |
| `linkbreakers qrcode-design public-create` | `linkbreakers qrcode-designs create` |
| `linkbreakers qrcode-design update` | `linkbreakers qrcode-designs update` |
| `linkbreakers qrcode-templates create-qrcode-template` | `linkbreakers qrcode-templates create` |
| `linkbreakers qrcode-templates delete-qrcode-template` | `linkbreakers qrcode-templates delete` |
| `linkbreakers qrcode-templates get-qrcode-template` | `linkbreakers qrcode-templates get` |
| `linkbreakers qrcode-templates import-qrcode-template` | `linkbreakers qrcode-templates import` |
| `linkbreakers qrcode-templates list-qrcode-templates` | `linkbreakers qrcode-templates list` |
| `linkbreakers qrcode-templates public-get-qrcode-template` | `linkbreakers qrcode-templates get-public` |
| `linkbreakers qrcode-templates public-list-qrcode-templates` | `linkbreakers qrcode-templates list-public` |
| `linkbreakers qrcode-templates publish-qrcode-template` | `linkbreakers qrcode-templates publish` |
| `linkbreakers qrcode-templates update-qrcode-template` | `linkbreakers qrcode-templates update` |
| `linkbreakers settings get-data-collection` | `linkbreakers link-settings get-data-collection` |
| `linkbreakers settings get-link-setting` | `linkbreakers link-settings get` |
| `linkbreakers settings update-data-collection` | `linkbreakers link-settings update-data-collection` |
| `linkbreakers settings update-link-setting` | `linkbreakers link-settings update` |
| `linkbreakers tags delete` | `linkbreakers tags bulk-delete` |
| `linkbreakers workflow-steps bulk-insert-workflow-step` | `linkbreakers workflow-steps bulk-insert` |
| `linkbreakers workflow-steps connect-step` | `linkbreakers workflow-steps connect` |
| `linkbreakers workflow-steps disconnect-step` | `linkbreakers workflow-steps disconnect` |
| `linkbreakers workflow-steps resolve-step-collisions` | `linkbreakers workflow-steps resolve-collisions` |
| `linkbreakers workspaces get-my-workspace` | `linkbreakers workspaces get` |
| `linkbreakers workspaces public-get-workspace` | `linkbreakers workspaces get-public` |

Four paths changed with them (`POST /v1/visitors:identify`,
`POST /v1/qrcode-designs`, `GET /v1/plan-usage`,
`POST /v1/custom-domains/{id}:verify`). The API removes a former path once no
client has called it for 14 days; none had an outside caller, so a CLI
released before this one gets `404` on all four former paths.

### Moved command (2026-11-16)

Resending a member's verification email moved to the Members section, so
`linkbreakers auth request-email-verification` becomes
`linkbreakers members request-email-verification`, and the `auth` group goes
away. It calls `POST /v1/members/email-verification`; the former
`POST /v1/auth/email/verification/request` is removed (only the dashboard
called it). Run it when a command fails with `WORKSPACE_UNVERIFIED`: it sends
the verification email again, and verifying the address lifts the pause.

### Google Forms is one of the form tools (2026-11-16)

Google Forms is now the tool `FORM_TOOL_GOOGLE_FORMS` of the form integrations
(`/v1/form-integrations/{tool}`), and `/v1/google-forms-integrations` is removed
(nothing outside the Linkbreakers dashboard called it over the 14 days before).
Regenerating removes the `linkbreakers google-forms-integrations` group:

| Before | After |
|---|---|
| `linkbreakers google-forms-integrations get` | `linkbreakers form-integrations get FORM_TOOL_GOOGLE_FORMS` |
| `linkbreakers google-forms-integrations update` | `linkbreakers form-integrations update FORM_TOOL_GOOGLE_FORMS` |
| `linkbreakers google-forms-integrations delete` | `linkbreakers form-integrations disconnect FORM_TOOL_GOOGLE_FORMS` |
| `linkbreakers google-forms-integrations list-detected-forms` | `linkbreakers form-integrations list-forms FORM_TOOL_GOOGLE_FORMS` |
| `linkbreakers google-forms-integrations sync-now` | `linkbreakers form-integrations prepare-forms FORM_TOOL_GOOGLE_FORMS` (reads the new responses now) |
| `linkbreakers google-forms-integrations choose-form-link` | `linkbreakers form-integrations choose-form-link FORM_TOOL_GOOGLE_FORMS <formId>` (new) |

(The exact argument form follows the generator.) Connecting Google Forms is a
browser sign-in: `linkbreakers form-integrations get-o-auth-url
FORM_TOOL_GOOGLE_FORMS` gives the URL to open.
