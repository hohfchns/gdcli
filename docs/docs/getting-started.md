To get started with **gdcli (Godot Project CLI)**, follow the steps below to download, install, and set up the tool on your Windows system.

---

## 1. Downloading gdcli

The latest version of gdcli is available on GitHub. It's recommended to download the `gdcliSetup.exe` for a straightforward installation process.

**Steps:**

1. Visit the [latest release page of gdcli on GitHub](https://github.com/IgorBayerl/gdcli/releases/latest).

2. Under the "Assets" section, locate and download the `gdcliSetup.exe` file.

---

## 2. Installing gdcli

After downloading the `gdcliSetup.exe`, proceed with the installation:

1. Run the `gdcliSetup.exe` file.

2. Follow the on-screen instructions to complete the installation.

The installer will automatically add gdcli to your system's PATH, allowing you to use the `gdcli` command from any command prompt window.

---

## 3. Alternative: Using the Standalone Binary

If you prefer not to use the installer, you can download the standalone binary (`gdcli.exe`) and manually add it to your system's PATH.

**Steps:**

1. Download the `gdcli.exe` from the [latest release page](https://github.com/IgorBayerl/gdcli/releases/latest).

2. Place the `gdcli.exe` file in a directory of your choice, for example, `C:\gdcli`.

3. Add this directory to your system's PATH:

   - **Open System Properties:**

     - Press the **Windows key** and type `environment variables`.

     - Select **"Edit the system environment variables"**.

   - **Access Environment Variables:**

     - In the **System Properties** window, click on the **"Advanced"** tab.

     - Click the **"Environment Variables..."** button.

   - **Edit the PATH Variable:**

     - In the **Environment Variables** window, locate the **"Path"** variable under **System variables**.

     - Select it and click **"Edit..."**.

   - **Add New Path:**

     - In the **Edit Environment Variable** window, click **"New"**.

     - Enter the path to the directory where you placed `gdcli.exe` (e.g., `C:\gdcli`).

     - Click **"OK"** to close all windows.

4. To verify the setup, open a new Command Prompt window and type:

   ```
   gdcli version
   ```

   If the installation was successful, this command will display the installed version of gdcli.

---

By following these steps, you'll have gdcli installed and ready to assist you in initializing and managing your Godot projects efficiently. 