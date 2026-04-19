import { Theme } from "@radix-ui/themes";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import ReactDOM from "react-dom/client";
import { ReactQueryProvider } from "./features/layout/ReactQuery";
import { routeTree } from "./routeTree.gen";
import "./styles.css";


const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  scrollRestoration: true,
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const rootElement = document.getElementById("app")!;

if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <Theme appearance="dark" accentColor="gray" grayColor="slate">
      <ReactQueryProvider>
        <RouterProvider router={router} />
      </ReactQueryProvider>
    </Theme>,
  );
}
