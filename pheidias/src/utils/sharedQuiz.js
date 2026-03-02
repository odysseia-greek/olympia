export function updateQuizUrl(router, currentQuery, routeName, updates = {}) {
    const newQuery = { ...currentQuery, ...updates };

    const queryChanged = Object.keys(updates).some(
        key => String(currentQuery[key] || '') !== String(updates[key])
    );

    if (queryChanged) {
        router.replace({ name: routeName, query: newQuery });
    }
}

export function shuffle(array) {
    const arr = [...array]
    for (let i = arr.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1))
        ;[arr[i], arr[j]] = [arr[j], arr[i]]
    }
    return arr
}