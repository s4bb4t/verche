# Verche

Verche updates dependencies in your `go.mod` file to the latest permissible versions and outputs a new file named `verched_go.mod`.

---

### **How It Works**

1. **Reads `go.mod`**  
   Parses the file line by line to find package dependencies and their current versions. Example:
   ```
   github.com/example/package v1.2.3
   ```

2. **Fetches Latest Versions**  
   Sends a request to a repository API (`https://repository.rt.ru/gateway/artifacts/findArtifacts`) to retrieve all available versions. Compares them to select the most recent version with the status `"PERMITTED"`.

3. **Updates Versions**  
   If a newer version is found, it replaces the current version. For example:
   ```
   github.com/example/package v1.2.3 --> github.com/example/package v1.3.0
   ```

4. **Writes Output**  
   Generates a new file, `verched_go.mod`, with updated versions while preserving the structure of the original file.
