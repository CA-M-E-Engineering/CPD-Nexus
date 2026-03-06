# CPD Submission Mapping Reference

This document details exactly how attendance and project data in the local database is mapped to the final JSON payload pushed to the SGBuildex Pitstop API.

## Payload Structure & Data Sources

The JSON payload structure follows the standard SGBuildex specification. Each data element is derived from specific database tables and columns.

### Participants Array
The `participants` array defines the regulatory body receiving the data.

| JSON Field | Source Table / Column | Description |
| :--- | :--- | :--- |
| `id` | `pitstop_authorisations.regulator_id` | Unique identifier of the regulator. Synced from Pitstop Config. |
| `name` | `pitstop_authorisations.regulator_name` | Name of the regulator (e.g., "BCA", "HDB"). Synced from Pitstop. |
| `meta.data_ref_id`| `projects.project_reference_number` | The project's official reference number. |

### OnBehalfOf Array
The `on_behalf_of` array is populated if the submission is being made on behalf of another registered entity.

| JSON Field | Source Table / Column | Description |
| :--- | :--- | :--- |
| `id` | `pitstop_authorisations.on_behalf_of_id` | UEN of the entity you are submitting on behalf of. Defaults to empty/null if direct submission. |

### Payload Data (Manpower Utilization)
The core data elements for each attendance record.

| JSON Field | Source Table / Column | Description |
| :--- | :--- | :--- |
| `submission_entity` | `projects.submission_entity` | Type of submitter (1 = Onsite Builder, 2 = Offsite Fabricator). |
| `submission_month` | `attendance.submission_date` | YYYY-MM format derived from the time of submission attempt. |
| `project_reference_number` | `projects.project_reference_number` | Official project reference. |
| `project_title` | `projects.project_title` | Official project title. |
| `project_location_description`| `projects.project_location_description` | Address of the project site. |
| `main_contractor_company_name` | `projects.main_contractor_name` | Name of the site main contractor. |
| `main_contractor_company_unique_entity_number` | `projects.main_contractor_uen` | UEN of the site main contractor. |
| `person_id_no` | `workers.person_id_no` | Worker FIN or NRIC. |
| `person_id_and_work_pass_type`| `workers.person_id_and_work_pass_type` | Work pass abbreviation (e.g., WP). |
| `person_nationality` | `workers.person_nationality` | Country code of worker origin. |
| `person_trade` | `workers.person_trade` | Worker's specific trade code. |
| `person_employer_company_name` | `projects.worker_company_name` | Name of the worker's direct employer. |
| `person_employer_company_unique_entity_number` | `projects.worker_company_uen` | UEN of the worker's direct employer. |
| `person_employer_company_trade` | `projects.worker_company_trade` | Array of trades performed by the employer. |
| `person_employer_client_company_name` | `projects.worker_company_client_name` | Name of the employer's client (who hired them). |
| `person_employer_client_company_unique_entity_number` | `projects.worker_company_client_uen` | UEN of the employer's client. |
| `person_attendance_date` | `attendance.time_in` | Date of the attendance record (YYYY-MM-DD). |
| `person_attendance_details.time_in` | `attendance.time_in` | ISO8601 UTC timestamp of worker entry. |
| `person_attendance_details.time_out`| `attendance.time_out` | ISO8601 UTC timestamp of worker exit. |

## Important Notes on Configuration

1. **Pitstop Authorisations**: The `regulator_id`, `regulator_name`, and `on_behalf_of_id` are **not** manually entered. They are populated by pulling the active configuration from the SGBuildex API using the "Sync Configuration" feature.
2. **Project Bonding**: When a project is created or edited in CPD Nexus, a specific "Pitstop Config" is bound to the project (`projects.pitstop_auth_id`).
3. **Empty Participant ID Error**: If `regulator_id` is empty, it means the synced `pitstop_authorisations` table row bound to this project does not have a valid ID stored. This usually requires going to "System Settings" -> "Pitstop Config" to resync with valid SGBuildex credentials.
