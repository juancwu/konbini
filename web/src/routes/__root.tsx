import { Navbar } from "@/components/ui/navbar";
import { createRootRoute, Outlet } from "@tanstack/react-router";

export const Route = createRootRoute({
	component: () => (
		<>
			<Navbar />
			<Outlet />
		</>
	),
});
