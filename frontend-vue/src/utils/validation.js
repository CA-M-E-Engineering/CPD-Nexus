/**
 * SGBuildex Validation Rules
 */

export const REGEX = {
    PROJECT_REF: /^[AE]\d{4}-[A-Za-z0-9]{5}-\d{4}$/,
    UEN: /^(?:\d{8}[A-Z]|\d{4}\d{5}[A-Z]|[TSR]\d{2}[A-Z0-9]{2}\d{4}[A-Z])$/,
    NRIC_FIN: /^[STFGM]\d{7}[A-Z0-9]$/,
    HDB_CONTRACT: /^D\/\d{1,5}\/\d{2}$/,
    LTA_CONTRACT: /^[A-Za-z0-9 /()-]{1,20}$/,
    SUBMISSION_MONTH: /^(?:19[5-9][0-9]|20[0-9]{2}|2100)-(?:0[1-9]|1[0-2])$/
};

export const validateProjectRef = (val) => REGEX.PROJECT_REF.test(val);
export const validateUEN = (val) => REGEX.UEN.test(val);
export const validateNRICFIN = (val) => REGEX.NRIC_FIN.test(val);
export const validateHDBContract = (val) => val.length <= 10 && REGEX.HDB_CONTRACT.test(val);
export const validateLTAContract = (val) => val.length <= 20 && REGEX.LTA_CONTRACT.test(val);

export const validatePersonTrade = (val) => {
    const validTrades = [
        '1.1', '1.2', '1.3', '1.4', '1.5',
        '2.1', '2.2', '2.3', '2.4', '2.5', '2.6', '2.7', '2.8',
        '3.1', '3.2', '3.3', '3.4', '3.5', '3.6', '3.7', '3.8', '3.9', '3.10', '3.11',
        '4.1', '4.2', '4.3', '4.4', '4.5', '4.6'
    ];
    return validTrades.includes(val);
};

export const validateWorkPassType = (val) => {
    const enums = ['SP', 'SB', 'EP', 'SPASS', 'WP', 'ENTREPASS', 'LTVP'];
    return enums.includes(val);
};

export const sanitizeUEN = (val) => (val || '').trim().toUpperCase();
