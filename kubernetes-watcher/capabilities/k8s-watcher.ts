import { Capability, a, Log } from "pepr";
import * as fs from "fs";

export const KubernetesWatcher = new Capability({
  name: "kubernetes-watcher",
  description: "Reports when Kubernetes EndpointSlice or Service changes",
});

const { When } = KubernetesWatcher;

interface EventLog {
  resource: string;
  timestamp: string;
  typeOfChange: string;
}

const LOG_FILE_PATH = "/tmp/k8s-watcher.log";

function KEvent(
  filePath: string,
  resource: string,
  typeOfChange: string,
): void {
  const logEntry: EventLog = {
    resource,
    timestamp: new Date().toISOString(),
    typeOfChange,
  };

  const logString = `${logEntry.timestamp} - Resource: ${logEntry.resource}, Type of Change: ${logEntry.typeOfChange}\n`;

  Log.info(logString);

  fs.appendFile(filePath, logString, err => {
    if (err) {
      Log.error(`Failed to write to file: ${err.message}`);
    } else {
      Log.debug("Log entry added successfully.");
    }
  });
}

When(a.EndpointSlice)
  .IsCreatedOrUpdated()
  .InNamespace("default")
  .WithName("kubernetes")
  .Reconcile(async () => {
    await KEvent(LOG_FILE_PATH, "CREATED_OR_UPDATED", "EndpointSlice");
  });

When(a.EndpointSlice)
  .IsDeleted()
  .InNamespace("default")
  .WithName("kubernetes")
  .Reconcile(async () => {
    await KEvent(LOG_FILE_PATH, "DELETED", "EndpointSlice");
  });

When(a.Service)
  .IsCreatedOrUpdated()
  .InNamespace("default")
  .WithName("kubernetes")
  .Reconcile(async () => {
    await KEvent(LOG_FILE_PATH, "CREATED_OR_UPDATED", "Service");
  });

When(a.Service)
  .IsDeleted()
  .InNamespace("default")
  .WithName("kubernetes")
  .Reconcile(async () => {
    await KEvent(LOG_FILE_PATH, "DELETED", "Service");
  });
