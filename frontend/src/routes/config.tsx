import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/config")({
    component: RouteComponent,
});

function RouteComponent() {
    return <div>Hello "/config"!</div>;
}
