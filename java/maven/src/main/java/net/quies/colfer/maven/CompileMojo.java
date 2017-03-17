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
 * Generates source code for a language. The options are: C, Go,
 * Java and JavaScript.
 * @author Pascal S. de Kloe
 */
@Mojo(name="compile", defaultPhase=LifecyclePhase.GENERATE_SOURCES)
public class CompileMojo extends AbstractMojo {

/** The target language. */
@Parameter(defaultValue="Java", required=true)
String lang;

/**
 * The source files. Directories are scanned for
 * files with the colf extension.
 */
@Parameter(defaultValue="src/main/colfer", required=true)
File[] schemas;

/** Normalizes schemas on the fly. */
@Parameter
boolean formatSchemas;

/** Adds a package prefix. Use slash as a separator when nesting. */
@Parameter
String packagePrefix;

/**
 * Sets the default upper limit for serial byte sizes. The
 * expression is applied to the target language under the name
 * ColferSizeMax. (default "16 * 1024 * 1024")
 */
@Parameter
String sizeMax;

/**
 * Sets the default upper limit for the number of elements in a
 * list. The expression is applied to the target language under the
 * name ColferListMax. (default "64 * 1024")
 */
@Parameter
String listMax;

/** Use a specific destination base directory. */
@Parameter(defaultValue="${project.build.directory}/generated-sources/colfer", required=true)
File sourceTarget;

@Component
MavenProject project;


@Override
public void execute()
throws MojoExecutionException, MojoFailureException {
	Path colf = compiler();

	try {
		Process proc = launch(colf);

		Scanner stderr = new Scanner(proc.getErrorStream());
		while (stderr.hasNext()) getLog().error(stderr.nextLine());

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
	args.add("-b=" + sourceTarget);
	if (packagePrefix != null)
		args.add("-p=" + packagePrefix);
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

		// ensure execution permission
		Set<PosixFilePermission> perms = Files.getPosixFilePermissions(path);
		if (!perms.contains(PosixFilePermission.OWNER_EXECUTE)) {
			perms.add(PosixFilePermission.OWNER_EXECUTE);
			Files.setPosixFilePermissions(path, perms);
		}

		return path;
	} catch (Exception e) {
		getLog().error("compiler command installation", e);
		throw new MojoExecutionException(path.toString() + ": installation failed");
	}
}

}
