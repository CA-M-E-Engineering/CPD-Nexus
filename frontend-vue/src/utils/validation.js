/**
 * SGBuildex Validation Rules — mirrors backend sgbuildex_rules.go
 * Keep in sync with: internal/pkg/validation/sgbuildex_rules.go
 */

// Pre-compiled regex patterns for SGBuildex field validation.
export const REGEX = {
    PROJECT_REF: /^[AE]\d{4}-[A-Za-z0-9]{5}-\d{4}$/,
    UEN: /^(?:\d{8}[A-Z]|\d{9}[A-Z]|[TSR]\d{2}[A-Z0-9]{2}\d{4}[A-Z])$/,
    NRIC_FIN: /^[STFGM]\d{7}[A-Z0-9]$/,
    HDB_CONTRACT: /^D\/\d{1,5}\/\d{2}$/,
    LTA_CONTRACT: /^[A-Za-z0-9 /()-]{1,20}$/,
    SUBMISSION_MONTH: /^(?:19[5-9]\d|20\d{2}|2100)-(?:0[1-9]|1[0-2])$/
};

// Valid sets for enum-based fields.
const VALID_TRADES = new Set([
    '1.1', '1.2', '1.3', '1.4', '1.5',
    '2.1', '2.2', '2.3', '2.4', '2.5', '2.6', '2.7', '2.8',
    '3.1', '3.2', '3.3', '3.4', '3.5', '3.6', '3.7', '3.8', '3.9', '3.10', '3.11',
    '4.1', '4.2', '4.3', '4.4', '4.5', '4.6'
]);

const VALID_PASS_TYPES = new Set(['SP', 'SB', 'EP', 'SPASS', 'WP', 'ENTREPASS', 'LTVP']);

// --- Exported validators ---

/** Validates a BCA project reference number (e.g. A1234-AB123-2022). */
export const validateProjectRef = (val) => REGEX.PROJECT_REF.test(val);

/** Validates a Singapore Unique Entity Number (UEN). */
export const validateUEN = (val) => REGEX.UEN.test(val?.toUpperCase?.() ?? val);

/** Validates a Singapore NRIC or FIN number. */
export const validateNRICFIN = (val) => REGEX.NRIC_FIN.test(val?.toUpperCase?.() ?? val);

/** Validates an HDB contract number (D/NNNNN/YY, max 10 chars). */
export const validateHDBContract = (val) => val?.length <= 10 && REGEX.HDB_CONTRACT.test(val);

/** Validates an LTA contract number (max 20 alphanumeric chars). */
export const validateLTAContract = (val) => val?.length <= 20 && REGEX.LTA_CONTRACT.test(val);

/** Validates a BCA trade code (e.g. '2.5', '3.11'). */
export const validatePersonTrade = (val) => VALID_TRADES.has(val);

/** Validates a work pass type against the accepted enum. */
export const validateWorkPassType = (val) => VALID_PASS_TYPES.has(val);

/** Normalises a UEN by trimming and uppercasing. */
export const sanitizeUEN = (val) => (val ?? '').trim().toUpperCase();
