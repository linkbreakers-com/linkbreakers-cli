# Changelog

## Unreleased (next release, from the API's public document of 2026-11-15)

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
`POST /v1/custom-domains/{id}:verify`). The API keeps answering the former
paths until at least 2027-06-01, so a CLI released before this one keeps
working until then.

### Moved command (2026-11-16)

Resending a member's verification email moved to the Members section, so
`linkbreakers auth request-email-verification` becomes
`linkbreakers members request-email-verification`, and the `auth` group goes
away. It calls `POST /v1/members/email-verification`; the former
`POST /v1/auth/email/verification/request` keeps answering until at least
2027-06-01. Run it when a command fails with `WORKSPACE_UNVERIFIED`: it sends
the verification email again, and verifying the address lifts the pause.
