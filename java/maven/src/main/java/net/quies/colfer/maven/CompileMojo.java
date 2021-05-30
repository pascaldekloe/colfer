package net.quies.colfer.maven;

import org.apache.maven.plugins.annotations.Component;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugin.MojoFailureException;
import org.apache.maven.project.MavenProject;

import java.io.IOException;
import java.io.InputStream;
import java.io.File;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.attribute.PosixFilePermission;
import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;
import java.util.Set;


/**
 * Compile Colfer schemas.
 * @author Pascal S. de Kloe
 */
@Mojo(name="compile", defaultPhase=LifecyclePhase.GENERATE_SOURCES)
public class CompileMojo extends AbstractMojo {

/**
 * The output is source code for either C, Dart, Go, Java or JavaScript.
 */
@Parameter(defaultValue="Java", required=true)
String lang;

/**
 * For each operand that names a file of a type other than
 * directory, colf reads the content as schema input. For each
 * named directory, colf reads all files with a .colf extension
 * within that directory.
 */
@Parameter(defaultValue="src/main/colfer", required=true)
File[] schemas;

/**
 * Normalize the format of all schema input on the fly.
 */
@Parameter
boolean formatSchemas;

/**
 * Use a specific base directory for the generated code.
 */
@Parameter(defaultValue="${project.build.directory}/generated-sources/colfer", required=true)
File sourceTarget;

/**
 * Compile to a package prefix.
 */
@Parameter
String packagePrefix;

/**
 * Supply custom tags. See the TAGS section in the manual for details.
 */
@Parameter
File[] tagFiles;

/**
 * Make all generated classes extend a super class. Java only.
 */
@Parameter
String superClass;

/**
 * Make all generated classes implement interfaces. Java only.
 */
@Parameter
String[] interfaces;

/**
 * Insert a code snippet from a file. Java only.
 */
@Parameter
File snippetFile;

/**
 * Sets the default upper limit for serial byte sizes. The
 * expression is applied to the target language under the name
 * ColferSizeMax. (default "16 * 1024 * 1024")
 */
@Parameter
String sizeMax;

/**
 * Sets the default upper limit for the number of elements in a
 * list. The expression is applied to the target language under
 * the name ColferListMax. (default "64 * 1024")
 */
@Parameter
String listMax;

@Component
MavenProject project;


@Override
public void execute()
throws MojoExecutionException, MojoFailureException {
	Path colf = compiler();

	try {
		Process proc = launch(colf);

		Scanner stderr = new Scanner(proc.getErrorStream());
		while (stderr.hasNext()) getLog().info(stderr.nextLine());

		int exit = proc.waitFor();
		if (exit != 0) throw new MojoFailureException("colf: exit " + exit);

		project.addCompileSourceRoot(sourceTarget.toString());
	} catch (Exception e) {
		throw new MojoExecutionException("compiler command execution", e);
	}
}

Process launch(Path colf)
throws IOException {
	List<String> args = new ArrayList<>();
	args.add(colf.toString());
	if (getLog().isDebugEnabled())
		args.add("-v");
	if (formatSchemas)
		args.add("-f");
	args.add("-b=" + sourceTarget);
	if (packagePrefix != null)
		args.add("-p=" + packagePrefix);
	if (tagFiles != null) {
		String[] paths = new String[tagFiles.length];
		for (int i = 0; i < paths.length; i++)
			paths[i] = tagFiles[i].getPath();
		args.add("-t=" + String.join(",", paths));
	}
	if (superClass != null)
		args.add("-x=" + superClass);
	if (interfaces != null)
		args.add("-i=" + String.join(",", interfaces));
	if (snippetFile != null)
		args.add("-c=" + snippetFile.getPath());
	if (sizeMax != null)
		args.add("-s=" + sizeMax);
	if (listMax != null)
		args.add("-l=" + listMax);
	args.add(lang);
	for (File s : schemas) args.add(s.toString());

	getLog().info("compile command arguments: " + args);
	ProcessBuilder builder = new ProcessBuilder(args);
	builder.directory(project.getBasedir());
	return builder.start();
}

/** Installs the executable. */
Path compiler()
throws MojoExecutionException {
	String command = "colf";
	String resource;
	{
		String arch = System.getProperty("os.arch").toLowerCase();
		if ("x86_64".equals(arch)) arch = "amd64";
		if (! "amd64".equals(arch))
			throw new MojoExecutionException("unsupported hardware architecture: " + arch);

		String os = System.getProperty("os.name", "generic").toLowerCase();
		if (os.startsWith("mac") || os.startsWith("darwin")) {
			resource = "/" + arch + "/colf-darwin";
		} else if (os.startsWith("windows")) {
			resource = "/" + arch + "/colf.exe";
			command = "colf.exe";
		} else {
			resource = "/" + arch + "/colf-" + os;
		}
	}
	Path path = new File(project.getBuild().getDirectory(), command).toPath();

	if (Files.exists(path)) return path;

	// install resource to path
	InputStream stream = CompileMojo.class.getResourceAsStream(resource);
	if (stream == null)
		throw new MojoExecutionException(resource + ": no such resource - platform not supported");
	try {
		Files.createDirectories(path.getParent());
		Files.copy(stream, path);
		stream.close();
	} catch (Exception e) {
		getLog().error("compiler command installation", e);
		throw new MojoExecutionException(path.toString() + ": installation failed");
	}

	try {
		if (path.getFileSystem().supportedFileAttributeViews().contains("posix")) {
			// ensure execution permission
			Set<PosixFilePermission> perms = Files.getPosixFilePermissions(path);
			if (!perms.contains(PosixFilePermission.OWNER_EXECUTE)) {
				perms.add(PosixFilePermission.OWNER_EXECUTE);
				Files.setPosixFilePermissions(path, perms);
			}
		}
	} catch (Exception e) {
		getLog().warn("compiler executable permission", e);
	}

	return path;
}

}
