import { beforeAll, afterAll } from "vitest";
import { userEvent, type UserEvent } from "vitest/browser";

export let user: UserEvent;

export function UserSetup() {
    beforeAll(() => {
        user = userEvent.setup();
    });
    afterAll(async () => {
        user.cleanup();
    });
}
