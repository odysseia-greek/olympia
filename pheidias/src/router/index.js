import { createRouter, createWebHistory } from 'vue-router';

const routerOptions = [
    { path: "/", view: "HomePage" },
    {
        path: "/quiz",
        view: "SokratesHome",
        name: "SokratesHome",
    },
    {
        path: "/quiz/media",
        view: "SokratesMedia",
        name: "QuizMedia",
        query: ["theme", "segment", "extendedResults", "doneAfter"],
    },
    {
        path: "/quiz/multiplechoice",
        view: "SokratesMultipleChoice",
        name: "QuizMultipleChoice",
        query: ["theme", "extendedResults", "doneAfter"],
    },
    {
        path: "/quiz/authorbased",
        view: "SokratesAuthorBased",
        name: "QuizAuthorBased",
        query: ["theme", "segment", "doneAfter"],
    },
    {
        path: "/quiz/dialogue",
        view: "SokratesDialogue",
        name: "QuizDialogue",
        query: ["theme", "set"],
    },
    {
        path: "/quiz/grammar",
        view: "SokratesGrammar",
        name: "QuizGrammar",
        query: ["theme", "segment", "set", "extendedResults", "doneAfter"],
    },
    {
        path: "/quiz/journey",
        view: "SokratesJourney",
        name: "Journey",
    },
    {
        path: "/texts",
        view: "Herodotos",
        name: "Herodotos",
        query: ["author", "book", "reference"],
    },
    {
        path: "/dictionary",
        view: "Alexandros",
        name: "Alexandros",
        query: ["language", "mode", "word", "extended"],
    },
    {
        path: "/grammar",
        view: "Dionysios",
        name: "Dionysios",
        query: ["word"],
    },
    { path: "/:pathMatch(.*)*", view: "NotFound" },
];

const routes = routerOptions.map((route) => {
    if (route.children) {
        return {
            ...route,
            component: () => import(`../views/${route.view}.vue`),
            children: route.children.map((child) => ({
                ...child,
                component: () => import(`../views/${child.view}.vue`),
            })),
        };
    }

    return {
        ...route,
        component: () => import(`../views/${route.view}.vue`),
    };
});


const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes,
});

export default router;
