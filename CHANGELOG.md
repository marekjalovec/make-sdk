## v0.5.2 [2024-05-06]

_What's new?_

- Token scope suggestions in case an API call fails moved to constants, and fixed where a wrong scope was suggested.
- Enums introduced for various cryptic values in the API.
- Scenario logs now contain author of changes to the scenario.

## v0.5.1 [2023-08-25]

_What's new?_

- github.com/google/uuid dependency removed

## v0.5.0 [2023-08-11]

_What's new?_

- Retry mechanism for API calls

## v0.4.0 [2023-04-19]

_What's new?_

- Endpoints supported

  - GET: `hooks/<id>`
  - GET: `hooks`

_Breaking changes_
- `Scenario`.`IsLinked` renamed to `IsEnabled`

## v0.3.0 [2023-03-22]

_What's new?_

- Endpoints supported

  - GET: `functions/<id>`
  - GET: `functions`

_Breaking changes_

  - `Organization`.`Licence` updated to match current documentation
  - `DataStore`.`DatastructureId` renamed to `DataStructureId`
 
## v0.2.0 [2023-02-14]

_What's new?_

- Endpoints supported

  - GET: `dlqs/<id>`
  - GET: `dlqs`
  - GET: `scenarios/<id>/logs/<log-id>`
  - GET: `scenarios/<id>/logs`
  - GET: `scenarios/<id>`
  - GET: `scenarios`

## v0.1.0 [2023-01-24]

_What's new?_

- Endpoints supported

  - GET: `connections/<id>`
  - GET: `connections`
  - GET: `data-stores/<id>`
  - GET: `data-stores`
  - GET: `organizations/<id>/variables`
  - GET: `organizations/<id>`
  - GET: `organizations`
  - GET: `teams/<id>/variables`
  - GET: `teams/<id>`
  - GET: `teams`
  - GET: `users/<id>/user-organization-roles`
  - GET: `users/<id>/user-team-roles`
  - GET: `users/me/api-tokens`
  - GET: `users/roles`
  - GET: `users`
