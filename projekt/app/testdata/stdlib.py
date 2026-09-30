"""A standard könyvtár néhány moduljának használata."""

import json
import math
from collections import Counter

words = "alma korte alma szilva alma korte".split()
counts = Counter(words)

report = {
    "most_common": counts.most_common(1)[0][0],
    "unique": len(counts),
    "sqrt_2": round(math.sqrt(2), 4),
}

print(json.dumps(report, sort_keys=True))
