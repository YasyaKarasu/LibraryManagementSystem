import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import Welcome from "./components/Welcome.vue";
import Book from "./components/Book.vue";
import Card from "./components/Card.vue";
import Borrow from "./components/Borrow.vue";

const routes = [
    { path: "/", component: Welcome},
    { name: "book", path: "/book", component: Book },
    { name: "card", path: "/card", component: Card },
    { name: "borrow", path : "/borrow", component: Borrow}
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

const app = createApp(App);

app.use(router);

app.mount("#app");
