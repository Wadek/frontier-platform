# Scorecard V_rev0 (synthetic OWASP fixtures)

Policy: english/SECURITY_POLICY_OWASP.md  
Runner: internal/owasp  
Suite: testdata/owasp (10 positive + 10 negative)

| Metric | Value |
|--------|-------|
| Positive cases (recall) | 10/10 TP (FN=0) |
| Negative cases (precision) | 10/10 TN (FP=0) |
| F1 on this suite | **1.0** |

DVNA lab: not yet F1=1 against full guidebook (only pattern V). See CUSTOMER_DVNA_BASELINE.md.  
Next: label DVNA intentional vulns → grow V until in-scope F1=1 → harden tree to gate.passed.
