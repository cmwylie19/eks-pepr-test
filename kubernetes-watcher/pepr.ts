import { PeprModule } from "pepr";
import cfg from "./package.json";
import { KubernetesWatcher } from "./capabilities/k8s-watcher";

new PeprModule(cfg, [KubernetesWatcher]);
